import { Box, Button, Card, CardContent, Paper, Typography } from '@mui/material';
import { EventTicket } from '../lib/actions/actions';

type Props = {
    ticket: EventTicket;
};

export default function Ticket({ ticket }: Props) {
    return (
        <Card>
            <CardContent>
                <Typography>{ticket.category}</Typography>
                <Box sx={{ marginBlock: '0.75rem' }}>
                    <Typography>
                        {ticket.crm.first_name} {ticket.crm.last_name}
                    </Typography>
                    <Typography>{ticket.crm.email}</Typography>
                </Box>
                <Typography>Bestilling: {ticket.order_id}</Typography>
            </CardContent>
            <Button variant="contained" color="primary">
                Tildel meg
            </Button>
            <Button variant="contained" color="primary">
                Tildel andre
            </Button>
        </Card>
    );
}
