import React, { useEffect, useState } from 'react';
import { Link as RouterLink, useHistory } from 'react-router-dom';
import 'firebase/auth';
import 'firebase/firestore';
import {
  FormControl,
  FormLabel,
  Input,
  Button,
  Heading,
  Divider,
  Box,
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
} from '@chakra-ui/react';
import { Helmet } from 'react-helmet';
import { useAuth, UserData } from '../contexts/AuthProvider';
import ContainerCentered from '../components/ContainerCentered';

export default function Login() {
  const auth = useAuth();
  const history = useHistory();
  const [values, setValues] = useState<UserData>({
    email: '',
    password: '',
  });

  useEffect(() => {
    if (auth.authenticated) history.push('/');
  }, []);

  const handleChange = (event: any) => {
    event.persist();
    setValues((values) => ({
      ...values,
      [event.target.name]: event.target.value,
    }));
  };

  const handleSubmit = (event: any) => {
    event.preventDefault();

    auth.login!(values, history).catch((error) => {
      console.log(error.message);
    });
  };

  return (
    <>
      <Helmet>
        <title>Login</title>
        <meta property="og:title" content="Login" />
      </Helmet>
      <ContainerCentered>
        <Breadcrumb>
          <BreadcrumbItem as={Heading} pb={4} isCurrentPage>
            <BreadcrumbLink as={RouterLink} to="/login">
              Log In
            </BreadcrumbLink>
          </BreadcrumbItem>
        </Breadcrumb>
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
      </ContainerCentered>
    </>
  );
}
