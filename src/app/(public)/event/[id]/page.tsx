import MainEvent from './event';

type Props = {
    params: { id: string };
};

const EventPage = ({ params: { id } }: Props) => {
    return <MainEvent id={id} />;
};
export default EventPage;
