import CardBase from './CardBase';

const MyEvents = async () => {
	return (
		<CardBase
			href="/my-events"
			description="Trykk for og gå til mine arrangementer"
			img="/my-events.jpg"
			imgAlt="Mine arrangementer"
			title="Mine arrangementer"
		/>
	);
};

export default MyEvents;
