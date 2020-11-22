import styles from '../styles/Home.module.css';
import React, {useEffect, useState} from 'react';

export default function Home() {
  const [users, setUsers] = useState([])

  useEffect(async ()=>{
    const resp = await fetch('http://localhost:8080/users');
    const respJSON = await resp.json();
    setUsers(respJSON);
  }, [])

  return (
    <div className='index'>
      {users.map(user=>{
        return(
          <div className='user' key={user.ID}>
            {user.ID}, {user.First} {user.Last}
          </div>
        );
      })}
    </div>
  );
}
