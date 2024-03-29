import React from 'react';
import {BrowserRouter as Router, Route, Routes} from 'react-router-dom';

import HomePage from './pages/Home';
import CaseListPage from './pages/CaseList'
import ProjectHomePage from './pages/ProjectHome'
import TestPlanListPage from './pages/TestPlanList'
import TestPlanPage from './pages/TestPlan'
import RunCaseListPage from './pages/RunCaseList'
import MilestonePage from './pages/Milestone'
import LoginPage from './pages/Login'
import SettingsPage from './pages/Settings';
import SettingMembersPage from './pages/SettingMembers';
import SettingPasswordPage from './pages/SettingPassword';
import SettingLanguagePage from './pages/SettingLanguage';
import PrivateRoute from './PrivateRoute';
import PrivateAdminRoute from './PrivateAdminRoute';

const App = () => {
    return (
        <Router>
            <Routes>
                <Route path="/login" element={<LoginPage/>}/>
                <Route path="/" element={<PrivateRoute><HomePage/></PrivateRoute>}/>
                <Route path="/settings" element={<PrivateRoute><SettingsPage/></PrivateRoute>}/>
                <Route path="/settings/members" element={<PrivateAdminRoute><SettingMembersPage/></PrivateAdminRoute>}/>
                <Route path="/settings/password" element={<PrivateRoute><SettingPasswordPage/></PrivateRoute>}/>
                <Route path="/settings/language" element={<PrivateRoute><SettingLanguagePage/></PrivateRoute>}/>
                <Route path="/:project_code/cases" element={<PrivateRoute><CaseListPage/></PrivateRoute>}/>
                <Route path="/:project_code/plans/:test_plan_id"
                       element={<PrivateRoute><TestPlanPage/></PrivateRoute>}/>
                <Route path="/:project_code/runs/:test_run_id"
                       element={<PrivateRoute><RunCaseListPage/></PrivateRoute>}/>
                <Route path="/:project_code/plans" element={<PrivateRoute><TestPlanListPage/></PrivateRoute>}/>
                <Route path="/:project_code/milestones" element={<PrivateRoute><MilestonePage/></PrivateRoute>}/>
                <Route path="/:project_code" element={<PrivateRoute><ProjectHomePage/></PrivateRoute>}/>
            </Routes>
        </Router>
    );
};

export default App;
