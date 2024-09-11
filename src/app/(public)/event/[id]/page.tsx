import { Metadata } from 'next';
import MainEvent from './components/MainEvent';
import { getAllEvents, getEventById } from '$app/(public)/components/lib/serverAction';

type Props = {
    params: { id: string };
};

export async function generateMetadata({ params }: Props): Promise<Metadata> {
    const id = params.id;
    const event = await getEventById(id);

    return {
        title: event.title,
        description: event.shortDescription,
    };
}

const EventPage = async ({ params: { id } }: Props) => {
    const event = await getAllEvents();
    const eventIndex = event.findIndex((event) => event.id === id);
    const prevNavigationId = event[eventIndex - 1]?.id;
    const nextNavigationId = event[eventIndex + 1]?.id;
    return <MainEvent id={id} prevNavigationId={prevNavigationId} nextNavigationId={nextNavigationId} />;
};
export default EventPage;
