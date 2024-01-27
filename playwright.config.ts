import { PlaywrightTestConfig } from '@playwright/test';

const config: PlaywrightTestConfig = {
    use: {
        baseURL: 'http://localhost:3001',
    },

    projects: [
        {
            name: 'Chromium',
            use: { browserName: 'chromium' },
        },
        {
            name: 'Firefox',
            use: { browserName: 'firefox' },
        },
        {
            name: 'WebKit',
            use: { browserName: 'webkit' },
        },
    ],

    testDir: './tests',
    testMatch: '**/*.spec.ts',

    timeout: 30000,

    workers: 2,
};

export default config;
