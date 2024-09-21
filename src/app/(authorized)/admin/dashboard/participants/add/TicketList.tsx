'use client';
import { Box, Card, CardActions, FormControl, Input, InputAdornment, InputLabel } from '@mui/material';
import SearchIcon from '@mui/icons-material/Search';
import React, { useEffect, useState } from 'react';
import { EventTicket } from '$app/(authorized)/my-profile/my-tickets/components/lib/actions/actions';
import Fuse from 'fuse.js';
import TicketCard from './components/TicketCard';

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
            {filteredTickets?.map((ticket) => <TicketCard key={ticket.id} ticket={ticket} />)}
        </Box>
    );
};

export default TicketList;
