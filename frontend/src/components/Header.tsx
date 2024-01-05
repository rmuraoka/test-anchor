import React from "react";
import {
    Avatar, Box,
    Button,
    ButtonGroup,
    Flex,
    Heading,
    IconButton,
    Menu,
    MenuButton,
    MenuItem,
    MenuList
} from "@chakra-ui/react";
import {SettingsIcon} from "@chakra-ui/icons";
import {useLocation, useNavigate} from "react-router-dom";
import {FiHome} from "react-icons/fi";
import LanguageSwitcher from "./LanguageSwitcher";

interface HeaderProps {
    project_code: string | undefined;
    is_show_menu: boolean;
}

const Header: React.FC<HeaderProps> = ({project_code, is_show_menu}) => {
    const navigate = useNavigate();
    const location = useLocation();
    const user = JSON.parse(localStorage.getItem('user') || '{}');
    const isActive = (path: string) => {
        if (path == '') {
            return location.pathname === `/${project_code}`
        }
        return location.pathname === `/${project_code}/${path}`;
    };
    const handleLogout = () => {
        localStorage.removeItem('token');
        localStorage.removeItem('user');
        navigate('/login');
    };
    const navigateTo = (path: string) => {
        navigate(`/${project_code}/${path}`);
    };
    const navigateToHome = () => {
        navigate('/');
    };
    const navigateToSetting = () => {
        navigate('/settings');
    };
    const navigateToUserAccounts = () => {
        navigate('/user/accounts');
    };
    return (
        <Flex
            as="header"
            align="center"
            justify="space-between"
            wrap="wrap"
            p="1rem"
            borderBottom="1px"
            borderColor="gray.200"
            position="fixed"
            top={0}
            left={0}
            right={0}
            zIndex="banner"
            bg="white"
            boxShadow="sm"
        >
            <Flex align="center" mr={5}>
                <Heading as="h1" size="lg" letterSpacing="tighter">
                    MyApp
                </Heading>
                <IconButton
                    aria-label="Home"
                    icon={<FiHome />}
                    colorScheme="gray"
                    onClick={navigateToHome}
                    ml={2}
                />
            </Flex>
            {is_show_menu && (
                <ButtonGroup variant="link" spacing={6}>
                    <Button
                        colorScheme={isActive('') ? 'blue' : 'gray'}
                        onClick={() => navigateTo('')}
                    >
                        Dashboard
                    </Button>
                    <Button
                        colorScheme={isActive('cases') ? 'blue' : 'gray'}
                        onClick={() => navigateTo('cases')}
                    >
                        Test Cases
                    </Button>
                    <Button
                        colorScheme={isActive('plans') ? 'blue' : 'gray'}
                        onClick={() => navigateTo('plans')}
                    >
                        Test Plans
                    </Button>
                    <Button
                        colorScheme={isActive('milestones') ? 'blue' : 'gray'}
                        onClick={() => navigateTo('milestones')}
                    >
                        Milestones
                    </Button>
                </ButtonGroup>
            )}
            <Flex align="center">
                <LanguageSwitcher />
                <IconButton
                    icon={<SettingsIcon/>}
                    aria-label="Settings"
                    variant="ghost"
                    mr={4}
                    ml={4}
                    onClick={() => navigateToSetting()}
                />
                <Menu>
                    <MenuButton as={Box}>
                        <Flex align="center">
                            <Avatar name={user.name} size="sm" />
                        </Flex>
                    </MenuButton>
                    <MenuList>
                        <MenuItem onClick={() => navigateToUserAccounts()}>Account</MenuItem>
                        <MenuItem onClick={handleLogout}>Logout</MenuItem>
                    </MenuList>
                </Menu>
            </Flex>
        </Flex>
    );
};

export default Header;