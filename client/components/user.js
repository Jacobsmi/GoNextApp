const user = (props) =>{
    const deleteClick = async (id) =>{
        const resp = await fetch(`http://localhost:8080/delete`,{
            method:"POST",
            body:JSON.stringify({
                'ID':id
            })
        });
        const respJSON = await resp.json();
        if(respJSON.Success == true){
            props.deleteUser(id)
        }
    }

    return(
      <div className="user">
          {props.id}, {props.first} {props.last} <button onClick={()=>{deleteClick(props.id)}}>Delete</button>
      </div>
    );
}

export default user;