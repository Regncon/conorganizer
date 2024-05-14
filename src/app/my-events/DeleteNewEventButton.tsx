'use client';

import IconButton from '@mui/material/IconButton';
import DeleteForeverOutlinedIcon from '@mui/icons-material/DeleteForeverOutlined';
import { db, firebaseAuth } from '$lib/firebase/firebase';
import { deleteDoc, doc } from 'firebase/firestore';
type Props = {
    docId: string;
};
const DeleteNewEventButton = ({ docId }: Props) => {
    const handleDeleteClick = async () => {
        if (firebaseAuth.currentUser?.uid && docId) {
            const eventRef = doc(db, 'users', firebaseAuth.currentUser?.uid, 'my-events', docId);
            await deleteDoc(eventRef);
        }
    };
    return (
        <IconButton
            sx={{
                color: '#f95e5e',
                position: 'absolute',
                padding: '1rem',
                bottom: '0',
                right: '0',
            }}
            onClick={handleDeleteClick}
        >
            <DeleteForeverOutlinedIcon />
        </IconButton>
    );
};

export default DeleteNewEventButton;
