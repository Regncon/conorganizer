'use client';
import {
    Box,
    Button,
    Card,
    CardActions,
    CardContent,
    FormControl,
    Input,
    InputAdornment,
    InputLabel,
    Typography,
} from '@mui/material';
import SearchIcon from '@mui/icons-material/Search';
import React, { useEffect, useState } from 'react';
import {
    ConvertTicketIdToParticipant,
    EventTicket,
} from '$app/(authorized)/my-profile/my-tickets/components/lib/actions/actions';
import Fuse from 'fuse.js';

type Props = {
    tickets: EventTicket[];
};

const TicketList = ({ tickets }: Props) => {
    const [searchQuery, setSearchQuery] = useState<string>('');
    const [filteredTickets, setFilteredTickets] = useState<EventTicket[]>(tickets);

    const fuseOptions = {
        keys: ['crm.first_name', 'crm.last_name', 'crm.email', 'order_id'],
        threshold: 0.5,
    };

    useEffect(() => {
        const dinnerTicketId = 157059;
        const fuse = new Fuse(
            tickets.filter((ticket) => ticket.category_id !== dinnerTicketId),
            fuseOptions
        );

        const result = searchQuery ? fuse.search(searchQuery).map((res) => res.item) : tickets;
        setFilteredTickets(result);
    }, [searchQuery, tickets]);

    const handleConvertToParticipant = async (id: number) => {
        const result = await ConvertTicketIdToParticipant(id);
        console.log(result);
    };
    return (
        <Box
            sx={{
                display: 'grid',
                gridTemplateColumns: 'repeat(auto-fit,minmax(306px, 1fr))',
                gap: '1rem',
            }}
        >
            <Card>
                <CardActions>
                    <FormControl variant="standard">
                        <InputLabel htmlFor="input-with-icon-adornment">Søk etter bilett</InputLabel>
                        <Input
                            id="input-with-icon-adornment"
                            value={searchQuery}
                            onChange={(e) => setSearchQuery(e.target.value)}
                            endAdornment={
                                <InputAdornment position="start">
                                    <SearchIcon />
                                </InputAdornment>
                            }
                        />
                    </FormControl>
                </CardActions>
            </Card>
            {filteredTickets?.map((ticket) => (
                <Card key={ticket.id}>
                    <CardContent>
                        <Typography>Bilett: {ticket.order_id}</Typography>
                        <Typography>{ticket.category}</Typography>
                        <Typography>
                            {ticket.crm.first_name} {ticket.crm.last_name}
                        </Typography>
                        <Typography>{ticket.crm.email}</Typography>
                    </CardContent>
                    <CardActions>
                        <Button
                            onClick={() => handleConvertToParticipant(ticket.id)}
                            variant="contained"
                            color="primary"
                        >
                            Konverter bilett til deltager
                        </Button>
                    </CardActions>
                </Card>
            ))}
        </Box>
    );
};

export default TicketList;
