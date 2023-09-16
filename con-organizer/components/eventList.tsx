'use client';

import { useAllEvents } from '@/lib/hooks/UseAllEvents';
// import db from '../lib/firebase';
import { Box, Card } from '../lib/mui';
import EventHeader from './eventHeader';

const EventList = () => {
    const { events, loading } = useAllEvents();
    console.log(events, loading);

    // const collectionRef = collection(db, 'events');
    // const [conEvents, setconEvents] = useState([] as ConEvent[]);
    // const [loading, setLoading] = useState(false);

    // useEffect(() => {
    //     setLoading(true);
    //     const unsub = onSnapshot(collectionRef, (querySnapshot) => {
    //         const items = [] as ConEvent[];
    //         querySnapshot.forEach((doc) => {
    //             items.push(doc.data() as ConEvent);
    //             items[items.length - 1].id = doc.id;
    //         });
    //         setconEvents(items);
    //         setLoading(false);
    //     });
    //     return () => {
    //         unsub();
    //     };
    // }, []);

    return (
        <Box className="flex flex-row flex-wrap justify-center gap-4 mb-20">
            {loading ? <h1>Loading...</h1> : null}
            {events?.map((conEvent) => (
                <Card
                    key={conEvent.id}
                    onClick={() => {
                        window.location.assign(`/event/${conEvent.id}`);
                    }}
                >
                    <EventHeader conEvent={conEvent} />
                </Card>
            ))}
        </Box>
    );
};

export default EventList;
