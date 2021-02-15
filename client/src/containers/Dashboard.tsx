import React from 'react';
import {
  Button,
  HStack,
  LinkOverlay,
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  Heading,
} from '@chakra-ui/react';
import { Helmet } from 'react-helmet';
import { useHistory, Link as RouterLink } from 'react-router-dom';
import { GoSignOut, GoGear } from 'react-icons/go';
import Card from '../components/Card';
import ContainerCentered from '../components/ContainerCentered';
import { useAuth } from '../contexts/AuthProvider';

export default function Dashboard() {
  const auth = useAuth();
  const history = useHistory();
  const handleClick = () => {
    auth.logout!(history).catch((error) => {
      console.log(error.message);
    });
  };

  return (
    <>
      <Helmet>
        <title>Dashboard</title>
        <meta property="og:title" content="Dashboard" />
      </Helmet>
      <ContainerCentered>
        <Breadcrumb>
          <BreadcrumbItem as={Heading} pb={4} isCurrentPage>
            <BreadcrumbLink as={RouterLink} to="/">
              Dashboard
            </BreadcrumbLink>
          </BreadcrumbItem>
        </Breadcrumb>

        <Card
          url="https://acolliercdn.ngrok.io/ABMrLUe.png"
          name="ABMrLUe.png"
        />
        <HStack spacing={4} mt="4">
          <Button
            leftIcon={<GoSignOut />}
            colorScheme="orange"
            onClick={handleClick}
          >
            Sign Out
          </Button>
          <Button leftIcon={<GoGear />} colorScheme="blue">
            <LinkOverlay as={RouterLink} to="/profile">
              Profile
            </LinkOverlay>
          </Button>
        </HStack>
      </ContainerCentered>
    </>
  );
}
