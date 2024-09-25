import CardBase from './CardBase';

const MyTickets = () => {
    return (
        <CardBase
            href="/my-profile/my-tickets"
            subTitle="Trykk for og gå til mine billetter"
            img="/my-tickets.webp"
            imgAlt="Mine billeter"
            title="Mine billetter"
        />
    );
};

export default MyTickets;
