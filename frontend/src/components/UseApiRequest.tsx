// useApiRequest.ts
import { useNavigate } from 'react-router-dom';

const API_URL = process.env.REACT_APP_BACKEND_URL;

export function useApiRequest() {
    const navigate = useNavigate();

    async function apiRequest(endpoint: string, options: RequestInit = {}) {
        const headers = new Headers({
            "Authorization": `Bearer ${localStorage.getItem('token')}`,
            'Content-Type': 'application/json',
        });

        const response = await fetch(`${API_URL}${endpoint}`, { headers, ...options });

        if (response.status === 401) {
            // トークンをクリア
            localStorage.removeItem('token');
            // ログインページにリダイレクト
            navigate('/login');
            return response;
        }

        if (!response.ok) {
            throw new Error(`API request failed: ${response.statusText}`);
        }

        return response;
    }

    return apiRequest;
}
