import React, { useEffect, useState } from 'react';
import {
  Button,
  Modal,
  ModalBody,
  ModalCloseButton,
  ModalContent,
  ModalFooter,
  ModalHeader,
  ModalOverlay,
  Link,
  useDisclosure,
  LinkOverlay,
} from '@chakra-ui/react';
import ShareXConfig from './ShareXConfig';
import { useAuth } from '../contexts/AuthContext';
import { useToken } from '../contexts/TokenContext';

export default function ShareXModal() {
  const { isOpen, onOpen, onClose } = useDisclosure();
  const auth = useAuth();
  const { token } = useToken(auth.user!.uid);
  const [config, setConfig] = useState('');
  const [configLink, setConfigLink] = useState('');

  useEffect(() => {
    const configStr = `{
      "Version": "13.4.0",
      "Name": "Custom CDN",
      "DestinationType": "ImageUploader",
      "RequestMethod": "POST",
      "RequestURL": "https://acolliercdn.ngrok.io/api/upload",
      "Headers": {
        "User-Agent": "ShareX",
        "Access-Token": "${token}"
      },
      "Body": "MultipartFormData",
      "FileFormName": "image",
      "RegexList": [
        "/https://acolliercdn.ngrok.io/upload/"
      ],
      "URL": "$json:url$"
    }
  `;

    setConfig(JSON.stringify(JSON.parse(configStr), null, 2));
  }, [token]);

  useEffect(() => {
    createConfigFile();
  }, [config]);

  const createConfigFile = () => {
    const data = new Blob([config], {
      type: 'application/vnd.sun.xml.calc',
    });
    if (configLink !== '') window.URL.revokeObjectURL(configLink);
    setConfigLink(window.URL.createObjectURL(data));
  };

  return (
    <>
      <Button h="1.75rem" size="sm" onClick={onOpen}>
        Open
      </Button>

      <Modal
        closeOnOverlayClick={false}
        isOpen={isOpen}
        onClose={onClose}
        isCentered
      >
        <ModalOverlay />
        <ModalContent>
          <ModalHeader>ShareX Config Download</ModalHeader>
          <ModalCloseButton />
          <ModalBody>
            <ShareXConfig config={config} />
          </ModalBody>

          <ModalFooter>
            <Button colorScheme="blue" mr={3} onClick={onClose}>
              Close
            </Button>
            <Button
              as={LinkOverlay}
              download="config.sxcu"
              href={configLink}
              textDecoration="none"
            >
              Download
            </Button>
          </ModalFooter>
        </ModalContent>
      </Modal>
    </>
  );
}
