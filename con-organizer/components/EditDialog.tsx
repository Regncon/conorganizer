'use client';

import { ErrorBoundary } from 'react-error-boundary';
import CloseIcon from '@mui/icons-material/Close';
import { Box, Dialog, Divider, IconButton } from '@mui/material';
import { ConEvent } from '@/models/types';
import EventBoundary from './ErrorBoundaries/EventBoundary';
import EditUi from './EditUi';
import EventUi from './EventUi';

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
