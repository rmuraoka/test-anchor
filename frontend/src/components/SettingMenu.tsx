import React from "react";
import {Box, Flex, Icon, Link, Text, VStack} from "@chakra-ui/react";
import {FaUser} from "react-icons/fa";

const Header: React.FC = () => {
    return (
        <Box w="20%" bg="gray.100" p={4} h="100vh" pt="6rem">
            <VStack align="stretch" spacing={4}>
                <Link href="/settings/members">
                    <Flex align="center">
                        <Icon as={FaUser} mr={2}/>
                        <Text fontSize="lg" fontWeight="bold">Member</Text>
                    </Flex>
                </Link>
            </VStack>
        </Box>
    );
};

export default Header;
