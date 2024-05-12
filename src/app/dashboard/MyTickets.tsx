import CardBase from './CardBase';

const MyTickets = () => {
    return (
        <CardBase
            href="/my-tickets"
            subTitle="Trykk for og gå til mine billetter"
            img="/my-tickets.jpg"
            imgAlt="Mine billeter"
            title="Mine billetter"
        />
    );
};

export default MyTickets;
