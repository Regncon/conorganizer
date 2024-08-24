import { getEventById } from '$app/(public)/serverAction';
import { Metadata } from 'next';
import MainEvent from './event';

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
