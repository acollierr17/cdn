import React from 'react';
import { Container } from '@chakra-ui/react';
import { Helmet } from 'react-helmet';
import Card from '../components/Card';

export default function Dashboard() {
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
      </Container>
    </>
  );
}
