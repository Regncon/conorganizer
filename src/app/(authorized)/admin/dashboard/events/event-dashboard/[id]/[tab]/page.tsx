import { notFound } from 'next/navigation';
import Edit from './edit/Edit';
import type { TabNames } from './lib/types/types';
import Room from './room/Room';
import Settings from './settings/Settings';
import type { Metadata } from 'next';
import { getEventById } from '$app/(public)/components/lib/serverAction';
import InterestPage from './interest/InterestPage';
import Players from './players/Players';

type Props = {
    params: {
        id: string;
        tab: TabNames;
    };
};
export async function generateMetadata({ params: { id, tab } }: Props): Promise<Metadata> {
    const event = await getEventById(id);

    let tabName;
    switch (tab) {
        case 'edit':
            tabName = 'Rediger';
            break;
        case 'interest':
            tabName = 'Ønsker';
            break;
        case 'players':
            tabName = 'Spillere';
            break;
        case 'room':
            tabName = 'Rom';
            break;
        case 'settings':
            tabName = 'Innstillinger';
            break;
        default:
            tabName = 'Ukjent';
            break;
    }

    return {
        title: `${tabName} | ${event.title}`,
    };
}

const page = ({ params: { id, tab } }: Props) => {
    if (tab === 'edit') {
        return <Edit id={id} />;
    }
    if (tab === 'room') {
        return <Room id={id} />;
    }
    if (tab === 'interest') {
        return <InterestPage id={id} />;
    }
    if (tab === 'players') {
        return <Players id={id} />;
    }
    if (tab === 'settings') {
        return <Settings id={id} />;
    }
    return notFound();
};

export default page;
