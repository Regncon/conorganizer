import CardBase from './CardBase';
const MyEvents = async () => {
    return (
        <CardBase
            href="/my-events"
            subTitle="Trykk for å gå til mine arrangement"
            img="/my-events.jpg"
            imgAlt="Mine arrangement"
            title="Mine arrangement"
        />
    );
};

export default MyEvents;
