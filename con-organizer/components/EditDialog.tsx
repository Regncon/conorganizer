'use client';

import { useEffect, useState } from 'react';
import { ErrorBoundary } from 'react-error-boundary';
import CloseIcon from '@mui/icons-material/Close';
import {
    Alert,
    Box,
    Dialog,
    DialogActions,
    DialogContent,
    DialogTitle,
    Divider,
    IconButton,
    InputLabel,
    MenuItem,
    Select,
    Switch,
    TextField,
} from '@mui/material';
import { doc, serverTimestamp, setDoc, updateDoc } from 'firebase/firestore';
import { eventsRef } from '@/lib/firebase';
import { GameType, Pool } from '@/models/enums';
import { ConEvent } from '@/models/types';
import { Button } from '../lib/mui';
import EventBoundary from './ErrorBoundaries/EventBoundary';
import EventUi from './EventUi';
import EditUi from './EditUi';

type Props = {
    open: boolean;
    conEvent?: ConEvent;
    handleClose: () => void;
};

const EditDialog = ({ open, conEvent, handleClose }: Props) => {
    return (
        <Dialog open={open} fullWidth={true} maxWidth="lg">
            <Box sx={{ minHeight: '900px', display: 'flex' }} flexDirection="row">
                <Box
                    className="p-4"
                    sx={{
                        width: '50%',
                        display: { xs: 'none', md: 'block' },
                    }}
                >
                    <ErrorBoundary FallbackComponent={EventBoundary}>
                        <EventUi conEvent={conEvent || ({} as ConEvent)} />
                    </ErrorBoundary>
                </Box>
                <Divider orientation="vertical" variant="middle" flexItem />

                <ErrorBoundary FallbackComponent={EventBoundary}>
                    <EditUi conEvent={conEvent || ({} as ConEvent)} />
                </ErrorBoundary>
                <Box sx={{ position: 'absolute', top: 0, right: 0 }}>
                    <IconButton onClick={handleClose} aria-label="close">
                        <CloseIcon />
                    </IconButton>
                </Box>
            </Box>
        </Dialog>
    );
};

export default EditDialog;
