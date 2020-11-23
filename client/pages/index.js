import styles from '../styles/Home.module.css';
import React, {useEffect, useState} from 'react';
import User from '../components/user';


export default function Home() {
  const [users, setUsers] = useState([])

  useEffect(async ()=>{
    const resp = await fetch('http://localhost:8080/users');
    const respJSON = await resp.json();
    setUsers(respJSON);
  }, [])

  const deleteUser = (id) =>{
    const newUsers = users.filter((user)=>{
      return user.ID !== id;
    })
    setUsers(newUsers);
  }

  return (
    <div className='index'>
      {users.map(user=>{
        return(
            <User key={user.ID} id={user.ID} first={user.First} last={user.Last} deleteUser={deleteUser} />
        );
      })}
    </div>
  );
}