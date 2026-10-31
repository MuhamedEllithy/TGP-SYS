"""Credential-handling safety tests; no infrastructure success is simulated."""
from contextlib import contextmanager, redirect_stdout
import io
from pathlib import Path
import stat
import subprocess
import tempfile
import unittest
from unittest.mock import patch
from urllib.parse import unquote, urlsplit

import local_env


class EnvironmentTests(unittest.TestCase):
    def setUp(self):
        self.directory = tempfile.TemporaryDirectory()
        self.addCleanup(self.directory.cleanup)
        self.root = Path(self.directory.name)
        self.path = self.root / '.env'
        self.root_patch = patch.object(local_env, 'ROOT', self.root)
        self.env_patch = patch.object(local_env, 'ENV_FILE', self.path)
        self.root_patch.start()
        self.env_patch.start()
        self.addCleanup(self.root_patch.stop)
        self.addCleanup(self.env_patch.stop)
        # Git checks are unit boundaries; real workspace exclusion is checked separately.
        def git(cmd, **kwargs):
            return subprocess.CompletedProcess(cmd, 1 if 'ls-files' in cmd else 0, b'', b'')
        self.git = patch.object(local_env.subprocess, 'run', side_effect=git)
        self.git.start()
        self.addCleanup(self.git.stop)
        self.socket = patch.object(local_env.socket, 'socket', side_effect=PermissionError)
        self.socket.start()
        self.addCleanup(self.socket.stop)
        self.output = io.StringIO()

    def generate(self):
        with redirect_stdout(self.output):
            local_env.configure()
        return local_env.load(self.path)

    def test_secure_generation_and_idempotence(self):
        env = self.generate()
        self.assertEqual(stat.S_IMODE(self.path.stat().st_mode), 0o600)
        secrets = [env[key] for key in ('DB_PASSWORD', 'POSTGRES_PASSWORD', 'MIGRATION_PASSWORD')]
        self.assertEqual(len(set(secrets)), 3)
        self.assertTrue(all(len(value) == 64 for value in secrets))
        self.assertTrue(all(value not in self.output.getvalue() for value in secrets))
        self.path.write_text(self.path.read_text() + '# unrelated setting\nEXTRA_VALUE=preserve-me\n')
        before = self.path.read_text()
        self.generate()
        self.assertEqual(self.path.read_text(), before)

    def test_new_host_configuration_skips_an_occupied_port(self):
        class Socket:
            def bind(self, address):
                if address[1] == 55432:
                    raise OSError("occupied")
        @contextmanager
        def socket():
            yield Socket()
        with patch.object(local_env.socket, 'socket', side_effect=socket):
            env = self.generate()
        self.assertEqual(env['DB_PORT'], '55433')
        self.assertEqual(urlsplit(env['DATABASE_URL']).port, 55433)
        self.assertEqual(urlsplit(env['TEST_DATABASE_URL']).port, 55433)

    def test_unrelated_environment_values_survive_initial_setup(self):
        self.path.write_text('# retained comment\nEXTRA_VALUE=preserve-me\n')
        env = self.generate()
        self.assertEqual(env['EXTRA_VALUE'], 'preserve-me')
        self.assertTrue(self.path.read_text().startswith('# retained comment\nEXTRA_VALUE=preserve-me\n'))
        self.assertIn('DATABASE_URL', env)

    def test_tracked_environment_is_never_written(self):
        with patch.object(local_env.subprocess, 'run', return_value=subprocess.CompletedProcess([], 0)):
            with self.assertRaises(ValueError):
                self.generate()
        self.assertFalse(self.path.exists())

    def test_inconsistent_existing_url_preserves_file(self):
        self.generate()
        text = '\n'.join(line.replace('tgp_app_dev:', 'other:') if line.startswith('DATABASE_URL=') else line for line in self.path.read_text().splitlines()) + '\n'
        self.path.write_text(text)
        with self.assertRaises(ValueError):
            self.generate()
        self.assertEqual(self.path.read_text(), text)

    def test_special_password_url_encoding(self):
        password = "space @:/?#%+$' and unicode \u0628"
        parsed = urlsplit(local_env.url('user', password, '127.0.0.1', 55432, 'db'))
        self.assertEqual(unquote(parsed.password), password)
        self.assertEqual(parsed.hostname, '127.0.0.1')
        self.assertEqual(parsed.path, '/db')

    def test_unrelated_existing_database_is_not_repointed(self):
        self.path.write_text('DATABASE_URL=unrelated\nOTHER=preserved\n')
        with self.assertRaises(ValueError):
            self.generate()
        self.assertEqual(self.path.read_text(), 'DATABASE_URL=unrelated\nOTHER=preserved\n')


if __name__ == '__main__':
    unittest.main()
