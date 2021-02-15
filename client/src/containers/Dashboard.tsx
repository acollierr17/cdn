import React from 'react';
import { Box, Button, Container, Text } from '@chakra-ui/react';
import { Helmet } from 'react-helmet';
import { useHistory } from 'react-router-dom';
import Card from '../components/Card';
import firebase from '../firebase';
import { useAuth } from '../contexts/AuthProvider';

export default function Dashboard() {
  const history = useHistory();
  const { user: currentUser } = useAuth();
  const handleClick = (event: any) => {
    event.preventDefault();

    firebase
      .auth()
      .signOut()
      .then(() => history.push('/login'));
  };

  return (
    <>
      <Helmet>
        <title>Dashboard</title>
        <meta property="og:title" content="Dashboard" />
      </Helmet>
      <Container
        pos="fixed"
        top="50%"
        left="50%"
        transform="translate(-50%, -50%)"
      >
        <Card
          url="https://acolliercdn.ngrok.io/ABMrLUe.png"
          name="ABMrLUe.png"
        />

        <Box>
          <Text>
            <strong>Email:</strong> {currentUser?.email ?? 'N/A'}
          </Text>
          <Button colorScheme="orange" mt="4" onClick={handleClick}>
            Sign Out
          </Button>
        </Box>
      </Container>
    </>
  );
}
