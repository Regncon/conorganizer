'use client';
import { Box, Button, Dialog, DialogContent, DialogContentText, DialogTitle, Typography } from '@mui/material';
import { useState } from 'react';

export interface SimpleDialogProps {
    open: boolean;
    onClose: () => void;
}

function SimpleDialog(props: SimpleDialogProps) {
    const { onClose, open } = props;

    const handleClose = () => {
        onClose();
    };

    return (
        <Dialog onClose={handleClose} open={open}>
            <DialogTitle>Under construction </DialogTitle>
            <DialogContent>
                <DialogContentText>Denne fuksjonen er ikke ferdig enda.</DialogContentText>
            </DialogContent>
            <img src="/under-construction.gif" alt="Under construction" />
        </Dialog>
    );
}

const FormSubmissionsPreviewPage = () => {
    const [open, setOpen] = useState(false);

    const handleClickOpen = () => {
        setOpen(true);
    };

    const handleClose = () => {
        setOpen(false);
    };

    return (
        <Box>
            <Typography variant="h1">Form Submissions Preview Page</Typography>
            <Button variant="contained" color="primary" onClick={handleClickOpen}>
                Godkjenn og start redigering
            </Button>
            <SimpleDialog open={open} onClose={handleClose} />
        </Box>
    );
};
export default FormSubmissionsPreviewPage;
