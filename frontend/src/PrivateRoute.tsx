// PrivateRoute.tsx
import React from 'react';
import { Navigate } from 'react-router-dom';
import {isAuthenticated, isAuthenticatedForAdmin} from './IsAuthenticated';

interface PrivateRouteProps {
    children: React.ReactNode;
}

const PrivateRoute: React.FC<PrivateRouteProps> = ({ children }) => {
    return isAuthenticated() ? <>{children}</> : <Navigate to="/login" />;
};

export default PrivateRoute;
