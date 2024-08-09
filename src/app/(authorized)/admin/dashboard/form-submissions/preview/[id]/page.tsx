'use client';
import MainEvent from '$app/(public)/event/[id]/event';
import { ConEvent } from '$lib/types';
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

    const event: ConEvent = {
        id: '1',
        title: 'Dette er tittelen',
        system: 'DnD 5e',
        gameMaster: 'Ola Nordmann',
        shortDescription: 'Dette er en kort beskrivelse',
        description: 'Dette er en lang beskrivelse',
        icons: ['katt', 'hund', 'fugl', 'rollespill', 'nisse', 'visse', 'nisse2', 'nisse3', 'nisse4'],
        published: false,
        email: '',
        name: '',
        phone: '',
        gameType: '',
        participants: 0,
        unwantedFridayEvening: false,
        unwantedSaturdayMorning: false,
        unwantedSaturdayEvening: false,
        unwantedSundayMorning: false,
        moduleCompetition: false,
        childFriendly: false,
        possiblyEnglish: false,
        adultsOnly: false,
        volunteersPossible: false,
        lessThanThreeHours: false,
        moreThanSixHours: false,
        beginnerFriendly: false,
        additionalComments: '',
        createdAt: '',
        createdBy: '',
        updateAt: '',
        updatedBy: '',
        subTitle: '',
        puljeFridayEvening: false,
        puljeSaturdayMorning: false,
        puljeSaturdayEvening: false,
        puljeSundayMorning: false,
    };
    return (
        <Box sx={{ maxWidth: '375px', margin: 'auto' }}>
            <Typography variant="h1">Forhåndsvisning</Typography>
            <Button variant="contained" color="primary" onClick={handleClickOpen}>
                Godkjenn og start redigering
            </Button>
            <SimpleDialog open={open} onClose={handleClose} />
            <hr />
            <MainEvent eventData={event} />
        </Box>
    );
};
export default FormSubmissionsPreviewPage;
