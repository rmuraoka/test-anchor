import React, {useState} from 'react';
import {Box, Button, ChakraProvider, FormControl, FormLabel, Input, useToast,} from '@chakra-ui/react';
import {useNavigate} from 'react-router-dom';

const Login: React.FC = () => {
    const API_URL = process.env.REACT_APP_BACKEND_URL
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const toast = useToast();
    const navigate = useNavigate(); // useNavigate フックの使用

    const handleSubmit = async (e: { preventDefault: () => void; }) => {
        e.preventDefault();
        try {
            const response = await fetch(`${API_URL}/login`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({"email": email, "password": password}),
            });
            const data = await response.json();
            if (response.ok) {
                localStorage.setItem('token', data.token);
                localStorage.setItem('user', JSON.stringify(data.user));
                navigate('/');
            } else {
                toast({
                    title: 'ログインに失敗しました',
                    description: "",
                    status: 'error',
                    duration: 5000,
                    isClosable: true,
                });
            }
        } catch (error) {
            console.error('Error:', error);
        }
    };

    return (
        <ChakraProvider>
            <Box maxWidth="sm" mx="auto" mt={10}>
                <form onSubmit={handleSubmit}>
                    <FormControl id="email" isRequired>
                        <FormLabel>email</FormLabel>
                        <Input
                            type="email"
                            value={email}
                            onChange={(e) => setEmail(e.target.value)}
                        />
                    </FormControl>
                    <FormControl id="password" mt={4} isRequired>
                        <FormLabel>password</FormLabel>
                        <Input
                            type="password"
                            value={password}
                            onChange={(e) => setPassword(e.target.value)}
                        />
                    </FormControl>
                    <Button mt={4} colorScheme="blue" type="submit">
                        Login
                    </Button>
                </form>
            </Box>
        </ChakraProvider>
    );
}

export default Login;
