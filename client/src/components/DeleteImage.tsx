import React, { useRef, useState } from 'react';
import {
  MenuItem,
  AlertDialog,
  AlertDialogOverlay,
  AlertDialogContent,
  AlertDialogHeader,
  AlertDialogBody,
  AlertDialogFooter,
  Button,
} from '@chakra-ui/react';
import { deleteImage } from '../api';

interface DeleteImageProps {
  fileName: string;
  token: string;
}

export default function DeleteImage(props: DeleteImageProps) {
  const [isOpen, setisOpen] = useState(false);
  const [loading, setLoading] = useState(false);
  const onDialogClose = () => setisOpen((state) => !state);
  const cancelRef = useRef<any>();

  const handleImageDeletion = async () => {
    try {
      setLoading(true);
      await deleteImage(props.fileName, props.token);
      setLoading(false);
      setisOpen(false);
    } catch (error) {
      setLoading(false);
      console.log(error.message);
    }
  };

  return (
    <>
      <MenuItem onClick={onDialogClose}>Delete</MenuItem>

      <AlertDialog
        leastDestructiveRef={cancelRef}
        isOpen={isOpen}
        onClose={onDialogClose}
      >
        <AlertDialogOverlay>
          <AlertDialogContent>
            <AlertDialogHeader fontSize="lg" fontWeight="bold">
              Delete Image
            </AlertDialogHeader>

            <AlertDialogBody>
              Are you sure? You can't undo this action afterwards.
            </AlertDialogBody>

            <AlertDialogFooter>
              <Button
                ref={cancelRef}
                onClick={onDialogClose}
                isLoading={loading}
              >
                Cancel
              </Button>
              <Button
                colorScheme="red"
                onClick={handleImageDeletion}
                ml={3}
                isLoading={loading}
              >
                Delete
              </Button>
            </AlertDialogFooter>
          </AlertDialogContent>
        </AlertDialogOverlay>
      </AlertDialog>
    </>
  );
}
