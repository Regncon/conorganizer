import { Metadata } from 'next';
import MainEvent from './components/MainEvent';
import { getAllEvents, getEventById } from '$app/(public)/components/lib/serverAction';
import MainEventBig from './components/MainEventBig';
import BigMediaQueryWrapper from './components/ui/BigMediaQueryWrapper';
import SmallMediaQueryWrapper from './components/ui/SmallMediaQueryWrapper';

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
    const events = await getAllEvents();
    const eventIndex = events.findIndex((event) => event.id === id);
    const prevNavigationId = events[eventIndex - 1]?.id;
    const nextNavigationId = events[eventIndex + 1]?.id;
    const event = await getEventById(id);
    return (
        <>
            <SmallMediaQueryWrapper>
                <MainEvent id={id} prevNavigationId={prevNavigationId} nextNavigationId={nextNavigationId} />
            </SmallMediaQueryWrapper>

            <BigMediaQueryWrapper>
                <MainEventBig event={event} prevNavigationId={prevNavigationId} nextNavigationId={nextNavigationId} />
            </BigMediaQueryWrapper>
        </>
    );
};
export default EventPage;
