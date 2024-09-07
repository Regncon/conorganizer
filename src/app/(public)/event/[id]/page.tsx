import { Metadata } from 'next';
import MainEvent from './components/MainEvent';
import { getEventById } from '$app/(public)/components/lib/serverAction';

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

const EventPage = ({ params: { id } }: Props) => {
    return <MainEvent id={id} />;
};
export default EventPage;
