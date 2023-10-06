import { CrmJson } from '@/models/types';

export async function GetParticipantsFromCheckIn() {
    // const url = https://app.checkin.no/graphql?client_id=API_KEY&client_secret=API_SECRET
    const query = `{
    eventTickets(customer_id: 13446, id: 58182, onlyCompleted: true) {
      id
      category
      category_id
      crm {
        first_name
        last_name
        id
        email
      }
      order_id
    }
  }`;

    const res = await fetch(
        `https://app.checkin.no/graphql?client_id=${process.env.CHECKIN_KEY}&client_secret=${process.env.CHECKIN_SECRET}`,
        {
            method: 'POST',
            body: JSON.stringify({ query }),
            headers: {
                'Content-Type': 'application/json',
            },
        }
    );

    const queryResult: CrmJson | undefined = await res.json();

    return queryResult;
}
