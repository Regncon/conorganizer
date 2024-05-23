'use client';
import { db, firebaseAuth } from '$lib/firebase/firebase';
import { faTrash } from '@fortawesome/free-solid-svg-icons/faTrash';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { buttonBaseClasses, buttonClasses } from '@mui/material';
import Box from '@mui/material/Box';
import IconButton from '@mui/material/IconButton';
import Typography from '@mui/material/Typography';
import { deleteDoc, doc } from 'firebase/firestore';
import { useState, type ComponentProps } from 'react';

type Props = {
    docId: string | undefined;
};

const TrashButton = ({ docId }: Props) => {
    const [disableRipple, setDisableRipple] = useState<boolean>(false);
    if (docId) {
        const handleDeleteClick: ComponentProps<'button'>['onClick'] = async (e) => {
            e.preventDefault();
            e.stopPropagation();
            setDisableRipple(true);

            if (firebaseAuth.currentUser?.uid && docId) {
                const eventRef = doc(db, 'users', firebaseAuth.currentUser?.uid, 'my-events', docId);
                await deleteDoc(eventRef);
            }
        };
        return (
            <IconButton
                component="span"
                className={[disableRipple ? 'disable-ripple' : ''].join(' ')}
                color="primary"
                onClick={handleDeleteClick}
            >
                <Box component={FontAwesomeIcon} icon={faTrash} />
            </IconButton>
        );
    }
    return null;
};

export default TrashButton;
