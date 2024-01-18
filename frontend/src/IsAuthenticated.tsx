export const isAuthenticated = () => {
    const token = localStorage.getItem('token');
    return !!token;
};

export const isAuthenticatedForAdmin = () => {
    const token = localStorage.getItem('token');
    const userStr = localStorage.getItem('user');
    if (!token || !userStr) {
        return false;
    }

    try {
        const user = JSON.parse(userStr);
        return user.permissions && user.permissions.includes('admin');
    } catch (error) {
        console.error('Error parsing user data from localStorage', error);
        return false;
    }
};