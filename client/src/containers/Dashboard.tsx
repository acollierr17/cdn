import React, { useEffect, useState } from 'react';
import {
  Button,
  Stack,
  LinkOverlay,
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  Heading,
  GridItem,
  Box,
  Flex,
  Spacer,
} from '@chakra-ui/react';
import { Helmet } from 'react-helmet';
import { useHistory, Link as RouterLink } from 'react-router-dom';
import { GoSignOut, GoGear } from 'react-icons/go';
import Card from '../components/Card';
import { useAuth } from '../contexts/AuthProvider';
import type { ImageResults } from '../api';
import { getAllImages } from '../api';
import { useToken } from '../contexts/TokenContext';

export default function Dashboard() {
  const auth = useAuth();
  const history = useHistory();
  const tokenState = useToken(auth.user!.uid);
  const [token, setToken] = useState(tokenState.token);
  const [loading, setLoading] = useState(true);

  const [images, setImages] = useState<ImageResults>({
    images: [],
    length: 0,
  });

  const handleClick = () => {
    auth.logout!(history).catch((error) => {
      console.log(error.message);
    });
  };

  const getImages = () => {
    setLoading(true);
    getAllImages(auth, token)
      .then((res) => {
        setImages(res);
        setLoading(false);
      })
      .catch((error) => {
        console.log(error.message);
        setLoading(false);
      });
  };

  useEffect(() => {
    if (!token) return;
    getImages();
  }, [token]);

  useEffect(() => {
    setToken(tokenState.token);
  }, [tokenState]);

  return (
    <>
      <Helmet>
        <title>Dashboard</title>
        <meta property="og:title" content="Dashboard" />
      </Helmet>

      <Box>
        <Flex alignItems="baseline" mx="25px">
          <Breadcrumb>
            <BreadcrumbItem as={Heading} pb={4} isCurrentPage>
              <BreadcrumbLink as={RouterLink} to="/">
                Dashboard
              </BreadcrumbLink>
            </BreadcrumbItem>
          </Breadcrumb>
          <Spacer />
          <Stack spacing={4} mt="4" direction={['column', 'row']}>
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
          </Stack>
        </Flex>
      </Box>

      <Flex flexDirection="row" flexWrap="wrap">
        {images.images.slice(0, 7).map((image, key) => {
          return (
            <>
              <GridItem m={5} w="250px">
                <Card key={key} loading={loading} token={token} {...image} />
              </GridItem>
            </>
          );
        })}
      </Flex>
    </>
  );
}
