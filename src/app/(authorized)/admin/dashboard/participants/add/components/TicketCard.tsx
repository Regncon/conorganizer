'use client';
import { Alert, Box, Button, Card, CardActions, CardContent, CircularProgress, Typography } from '@mui/material';
import {
    ConvertTicketIdToParticipant,
    EventTicket,
} from '$app/(authorized)/my-profile/my-tickets/components/lib/actions/actions';
import { useState } from 'react';
import { ActionResponse } from '$lib/types';
import Over18 from '$ui/participant/Over18';

type Props = {
    ticket: EventTicket;
};

const TicketCard = ({ ticket }: Props) => {
    const [convertResponce, setConvertResponce] = useState<ActionResponse>();
    const [loading, setLoading] = useState(false);

    const handleConvertToParticipant = async (id: number) => {
        setLoading(true);
        const result = await ConvertTicketIdToParticipant(id);

        if (result.type === 'error') {
            console.error(result);
        }

        setConvertResponce(result);
        setLoading(false);
    };

    const isOver18 = new Date().getFullYear() - new Date(ticket.crm.born).getFullYear() > 18;

    return (
        <Card key={ticket.id}>
            <CardContent>
                <Typography>Bestilling: {ticket.order_id}</Typography>
                <Typography>{ticket.category}</Typography>
                <Typography>
                    {ticket.crm.first_name} {ticket.crm.last_name}
                </Typography>
                <Typography>{ticket.crm.email}</Typography>
                <Box sx={{ display: 'flex', justifyContent: 'start', alignItems: 'center' }}>
                    <Over18 over18={isOver18} />
                </Box>
                {loading ?
                    <CircularProgress />
                    : convertResponce && <Alert severity={convertResponce.type}>{convertResponce.message}</Alert>}
            </CardContent>
            <CardActions>
                <Button
                    onClick={() => handleConvertToParticipant(ticket.id)}
                    variant="contained"
                    color="primary"
                    disabled={loading}
                >
                    Konverter bilett til deltager
                </Button>
            </CardActions>
        </Card>
    );
};

export default TicketCard;
