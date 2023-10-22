import { Card } from '@mui/material';
import { Route } from 'next';
import Link from 'next/link';
import { ConEvent } from '@/models/types';
import EventCardHeader from '../EventList/EventCardHeader';

type Props = {
    conEvent: ConEvent;
};

const EventCard = ({ conEvent }: Props) => {
    return (
        <Card
            key={conEvent.id}
            component={Link}
            href={`/event/${conEvent.id}` as Route}
            sx={{
                width: { xs: '100vw', md: '500px' },
                cursor: 'pointer',
                opacity: conEvent?.published === false ? '50%' : '',
                textDecoration: 'none',
                display: 'grid',
                gridTemplateRows: '1fr auto',
            }}
        >
            <EventCardHeader conEvent={conEvent} />
        </Card>
    );
};

export default EventCard;
