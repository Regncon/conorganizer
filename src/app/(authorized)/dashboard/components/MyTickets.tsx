import CardBase from './CardBase';

const MyTickets = () => {
    // Retunr null to prevent errors in early access
    return null;
    return (
        <CardBase
            href="/my-profile/my-tickets"
            subTitle="Trykk for og gå til mine billetter"
            img="/my-tickets.jpg"
            imgAlt="Mine billeter"
            title="Mine billetter"
        />
    );
};

export default MyTickets;
