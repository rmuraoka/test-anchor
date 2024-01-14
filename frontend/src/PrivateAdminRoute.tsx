import React from 'react';
import {Navigate} from 'react-router-dom';
import {isAuthenticatedForAdmin} from './IsAuthenticated';

interface PrivateRouteProps {
    children: React.ReactNode;
}

const PrivateAdminRoute: React.FC<PrivateRouteProps> = ({children}) => {
    if (isAuthenticatedForAdmin()) {
        return <>{children}</>
    } else {
        localStorage.removeItem('token');
        localStorage.removeItem('user');
        return <Navigate to="/login"/>
    }
};

export default PrivateAdminRoute;
