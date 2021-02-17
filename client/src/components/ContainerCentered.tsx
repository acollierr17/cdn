import React from 'react';
import { Container } from '@chakra-ui/react';

export default function ContainerCentered({ children, ...rest }: any) {
  return (
    <Container
      pos="fixed"
      top="50%"
      left="50%"
      transform="translate(-50%, -50%)"
      {...rest}
    >
      {children}
    </Container>
  );
}
