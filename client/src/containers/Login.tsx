import React, { useState } from 'react';
import { useHistory } from 'react-router-dom';
import firebase from '../firebase';
import 'firebase/auth';
import 'firebase/firestore';
import {
  Container,
  FormControl,
  FormLabel,
  Input,
  Button,
  Heading,
  Divider,
  Box,
} from '@chakra-ui/react';
import { Helmet } from 'react-helmet';
import { useAuth } from '../contexts/AuthProvider';

interface UserData {
  email: string;
  password: string;
}

export default function Login() {
  const auth = useAuth();
  const history = useHistory();
  const [values, setValues] = useState<UserData>({
    email: '',
    password: '',
  });

  const handleChange = (event: any) => {
    event.persist();
    setValues((values) => ({
      ...values,
      [event.target.name]: event.target.value,
    }));
  };

  const handleSubmit = (event: any) => {
    event.preventDefault();
    console.log(import.meta.env);

    firebase
      .auth()
      .signInWithEmailAndPassword(values.email, values.password)
      .then((res) => {
        auth.setUser(res);
        history.push('/');
      })
      .catch((error) => {
        console.log(error.message);
      });
  };

  return (
    <>
      <Helmet>
        <title>Login</title>
        <meta property="og:title" content="Login" />
      </Helmet>
      <Container
        pos="fixed"
        top="50%"
        left="50%"
        transform="translate(-50%, -50%)"
      >
        <Heading pb="4">Log In</Heading>
        <Divider />
        <Box pt="4">
          <form onSubmit={handleSubmit}>
            <FormControl id="email">
              <FormLabel>Email address</FormLabel>
              <Input
                type="email"
                name="email"
                placeholder="Enter your email"
                onChange={handleChange}
                isRequired
              />
            </FormControl>

            <FormControl id="password">
              <FormLabel>Password</FormLabel>
              <Input
                type="password"
                name="password"
                placeholder="Enter your password"
                onChange={handleChange}
                isRequired
              />
            </FormControl>

            <Button colorScheme="blue" type="submit" mt="4">
              Log In
            </Button>
          </form>
        </Box>
      </Container>
    </>
  );
}
