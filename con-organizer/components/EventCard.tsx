import { Card } from '@mui/material';
import { Route } from 'next';
import Link from 'next/link';
import { ConEvent } from '@/models/types';
import EventHeader from './EventHeader';

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
                maxWidth: '500px',
                cursor: 'pointer',
                opacity: conEvent?.published === false ? '50%' : '',
                textDecoration: 'none',
            }}
        >
            <EventHeader conEvent={conEvent} listView={true} />
        </Card>
    );
};

export default EventCard;
