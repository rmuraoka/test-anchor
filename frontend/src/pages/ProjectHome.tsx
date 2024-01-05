import React, {useEffect, useState} from 'react';
import {Box, ChakraProvider, Flex, Heading, SimpleGrid, Text, useColorModeValue} from '@chakra-ui/react';
import {useParams} from "react-router-dom";
import Header from "../components/Header";
import {useApiRequest} from "../components/UseApiRequest";

interface User {
    id: number;
    name: string;
}

interface Milestone {
    id: number;
    title: string;
    description: string;
    due_date: string;
    status: string;
    test_case_count: number;
}

interface TestPlan {
    id: number;
    project_id: number;
    title: string;
    count: number;
    status: string;
    created_by: User;
    updated_by: User;
}

interface Project {
    id: number;
    code: string;
    title: string;
    description: string;
    milestones: Milestone[];
    test_plans: TestPlan[];
}

interface CardProps {
    title: string;
    children: React.ReactNode;
}

const ProjectHome: React.FC = () => {
    const [project, setProject] = useState<Project | null>(null);
    const {project_code} = useParams();
    const apiRequest = useApiRequest();

    // APIからテストケースを取得
    const fetchProjects = async () => {
        try {
            const response = await apiRequest(`/protected/projects/${project_code}`);
            const data = await response.json();
            setProject(data)
        } catch (error) {
            console.error('Error fetching TestCases:', error);
        }
    };

    const Card: React.FC<CardProps> = ({title, children}) => {
        const cardBg = useColorModeValue('white', 'gray.700');
        return (
            <Box
                p={5}
                shadow="md"
                borderWidth="1px"
                borderRadius="lg"
                bg={cardBg}
                mb={4}
            >
                <Heading fontSize="xl">{title}</Heading>
                <Text mt={4}>{children}</Text>
            </Box>
        );
    };

    useEffect(() => {
        fetchProjects();
    }, [project_code]);

    return (
        <ChakraProvider>
            <Flex direction="column">
                {/* ヘッダー */}
                <Header project_code={project_code} is_show_menu={true}/>
                {/* 進捗状況ダッシュボード */}
                <Box p={8} pt="6rem">
                    <Heading as="h3" size="lg" mb={6}>Progress Dashboard</Heading>
                    <Flex>
                        {/* テストプランセクション */}
                        <Box w="50%" pr={2}>
                            <Heading as="h4" size="md" mb={4}>Test Plans</Heading>
                            <SimpleGrid columns={1} spacing={4}>
                                {project?.test_plans.map(plan => (
                                    <Card key={plan.id} title={plan.title}>
                                        Status: {plan.status}
                                    </Card>
                                ))}
                            </SimpleGrid>
                        </Box>
                        {/* マイルストーンセクション */}
                        <Box w="50%" pl={2}>
                            <Heading as="h4" size="md" mb={4}>Milestones</Heading>
                            <SimpleGrid columns={1} spacing={4}>
                                {project?.milestones.map(milestone => (
                                    <Card key={milestone.id} title={milestone.title}>
                                        Due Date: {milestone.due_date}
                                    </Card>
                                ))}
                            </SimpleGrid>
                        </Box>
                    </Flex>
                </Box>
            </Flex>
        </ChakraProvider>
    );
};

export default ProjectHome;
