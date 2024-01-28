import { test, expect } from '@playwright/test';
require('dotenv').config({ path: '.env.test' });

test('Login and verify that test cases can be executed', async ({ page }) => {
    let testPassword = process.env.TEST_LOGIN_PASSWORD;

    // Access Mailhog if the password is not set in the environment variable
    if (testPassword === '') {
        await page.goto('http://localhost:8026');

        // Search for the email with the title "Test Anchor Account"
        await page.fill('input#search', 'Test Anchor Account');
        await page.keyboard.press('Enter');
        await page.waitForSelector('text=Test Anchor Account');

        // Select the targeted email
        await page.click('span.subject.unread:has-text("Test Anchor Account")');
        await page.waitForSelector('#preview-plain');

        // Retrieve the password from the email
        const emailBody = await page.textContent('#preview-plain');
        const passwordMatch = emailBody.match(/Your Password is (\S+)/);
        if (!passwordMatch) throw new Error('Password not found in email');
        testPassword = passwordMatch[1];
    }

    // Retrieve login information from environment variables
    const testEmail = process.env.INITIAL_USER_EMAIL;

    // Access the login page
    await page.goto('http://localhost:3001/login');

    // Enter login credentials
    await page.fill('input[type="email"]', testEmail);
    await page.fill('input[type="password"]', testPassword);

    // Click the login button
    await page.click('button:has-text("Login")');

    // Wait for the post-login redirect
    await page.waitForURL('http://localhost:3001/');

    // Click the button to add a new project
    await page.click('button:has-text("Add a new Project")');

    // Wait for the modal to appear
    await page.waitForSelector('text=Add a new project');

    // Enter project information
    const timestamp = Date.now();
    const uniqueProjectCode = `NP${timestamp}`;
    await page.fill('[data-test="input title"]', uniqueProjectCode);
    await page.fill('[data-test="input code"]', uniqueProjectCode);
    await page.fill('[data-test="input description"]', 'Description of new project');
    await page.keyboard.press('Tab');
    await page.waitForTimeout(300);

    // Add the project
    await page.click('[data-test="button project add"]');

    // Move to the project page
    await page.click(`tbody tr:has-text("${uniqueProjectCode}") td:nth-child(1) a`);

    // Navigate to the test cases page
    await page.click('button:has-text("Test Cases")');
    await page.screenshot({ path: 'screenshot.png' });

    // Add a test suite
    await page.click('button[aria-label="Add test suite"]');
    await page.fill('[data-test="input testsuite name"]', "Test Suite Name");
    await page.click('[data-test="button testsuite add"]');

    // Add a test case
    await page.click('button[aria-label="Options"]');
    await page.click('button:has-text("Add case")');
    await page.fill('[data-test="input testcase title"]', "Test Case Name");
    await page.fill('[data-test="textarea testcase content"]', "Test Case Content");
    await page.click('[data-test="button testcase add"]');
    await page.waitForTimeout(300);

    await page.screenshot({ path: 'screenshot.png' });
    await page.click('button[aria-label="Open Test Case"]');

    // Check that the test case contains the added content
    expect(await page.textContent('body')).toContain('Test Case Content');

    // Navigate to the test plans page
    await page.click('button:has-text("Test Plans")');

    // Add a test plan
    await page.click('button:has-text("Add New Test Plan")');
    await page.fill('[data-test="input testplan title"]', "Test Plan Title");
    await page.click('[data-test="button testplan add"]');

    // Move to the test plan
    await page.click(`tbody tr:has-text("Test Plan Title") td:nth-child(1) a`);

    // Add a test run
    await page.click('button:has-text("Add New Test Run")');
    await page.fill('[data-test="input testrun title"]', "Test Run Title");
    await page.click('[data-test="button testrun add"]');

    // Edit the test run
    await page.click('button:has-text("Edit")');
    await page.waitForTimeout(300);
    const firstCheckboxControl = await page.$('span.chakra-checkbox__control');
    if (firstCheckboxControl) {
        await firstCheckboxControl.click();
    }
    await page.click('button:has-text("Update")');

    // Move to the test run
    await page.click(`tbody tr:has-text("Test Run Title") td:nth-child(1) a`);

    // Change the status of the test case
    await page.click('button:has-text("Untested")');
    await page.click('button:has-text("Passed")');
    await page.waitForTimeout(300);

    // Verify that the status has been updated
    expect(await page.textContent('body')).toContain('Status has been updated.');
});
