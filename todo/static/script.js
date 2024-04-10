// Обработчик события DOMContentLoaded.
document.addEventListener('DOMContentLoaded', async function() {
  await refreshTaskList();
 });
 
 // Обработчик отправки формы создания задачи.
 document.getElementById('task-form').addEventListener('submit', async function(e) {
     e.preventDefault();
     const taskInput = document.getElementById('task-input');
     const taskText = taskInput.value.trim();
   
     if (taskText === '') {
       alert('Введите текст задачи');
       return;
     }

     if (taskText.length > 255) {
        alert('Текст задачи не может превышать 255 символов');
        return;
      }
   
     const currentDate = new Date();
     const task = {
       text: taskText,
       date: currentDate,
       status: 'в процессе'
     };
   
     try {
       await createTask(task);
       taskInput.value = '';
       await refreshTaskList();
     } catch (error) {
       console.error('Error when creating a task:', error);
       alert('An error occurred while creating a task. Please try again.');
     }
 });
 
 // Обработчик кликов по списку задач (удаление и редактирование).
 document.getElementById('task-list').addEventListener('click', async function(e) {
  if (e.target.classList.contains('delete-btn')) {
       const taskId = e.target.parentElement.dataset.taskId;
       try {
           await deleteTask(taskId);
           await refreshTaskList();
       } catch (error) {
           console.error('Error when deleting a task:', error);
       }
  }
 
  if (e.target.classList.contains('edit-btn')) {
     const taskItem = e.target.parentElement;
     const taskId = taskItem.dataset.taskId;
     const taskText = taskItem.querySelector('.task-text');
     const editInput = taskItem.querySelector('.edit-input');
     const statusSelect = taskItem.querySelector('.status-select');
 
     if (editInput.style.display === 'none') {
         editInput.style.display = 'inline';
         editInput.value = taskText.textContent;
         taskText.style.display = 'none';
         statusSelect.style.display = 'none';
         e.target.textContent = 'Сохранить';
     } else {
         const currentDate = new Date();
         const updatedTask = {
             id: taskId,
             text: editInput.value.trim(),
             status: statusSelect.value,
             date: currentDate
         };
 
         // Валидация полей задачи.
         if (updatedTask.text === '') {
             alert('Текст задачи не может быть пустым');
             return;
         }

         if (updatedTask.text.length > 255) {
            alert('Текст задачи не может превышать 255 символов');
            return;
          }
 
         if (updatedTask.status !== 'в процессе' && updatedTask.status !== 'завершено') {
             alert('Некорректный статус задачи');
             return;
         }
 
         try {
             await updateTask(updatedTask);
             taskText.textContent = editInput.value;
             taskItem.querySelector('.status-select').value = statusSelect.value;
             taskItem.querySelector('.task-date').textContent = currentDate.toLocaleString();
             editInput.style.display = 'none';
             taskText.style.display = 'inline';
             statusSelect.style.display = 'inline';
             e.target.textContent = 'Редактировать';
             await refreshTaskList();
         } catch (error) {
             console.error('Error when updating a task:', error);
         }
     }
  }
 });
 
 // Обработчик изменения фильтра по статусу.
 document.getElementById('status-filter').addEventListener('change', async function() {
  await refreshTaskList();
 });
 
 // Обработчик изменения фильтра сортировки.
 document.getElementById('sort-filter').addEventListener('change', async function() {
  await refreshTaskList();
 });
 
 // Обработчик изменения статуса задачи.
 document.getElementById('task-list').addEventListener('change', async function(e) {
  if (e.target.classList.contains('status-select')) {
       const taskItem = e.target.closest('.task-item');
       const taskId = taskItem.dataset.taskId;
       const taskText = taskItem.querySelector('.task-text').textContent;
       const statusSelect = taskItem.querySelector('.status-select');
       const currentDate = new Date();
 
       const updatedTask = {
           id: taskId,
           text: taskText,
           status: statusSelect.value,
           date: currentDate
       };
 
       try {
           await updateTask(updatedTask);
           await refreshTaskList();
         } catch (error) {
           console.error('Error when updating task status:', error);
         }
  }
 });
 
 // Функция создания новой задачи.
 async function createTask(task) {
  const response = await fetch('/api/tasks/create', {
       method: 'POST',
       headers: {
           'Content-Type': 'application/json'
       },
       body: JSON.stringify(task)
  });
 
  if (!response.ok) {
       throw new Error('Error when creating a task');
  }
 }
 
 // Функция обновления задачи.
 async function updateTask(task) {
  const response = await fetch(`/api/tasks/update?id=${task.id}`, {
       method: 'PUT',
       headers: {
           'Content-Type': 'application/json'
       },
       body: JSON.stringify(task)
  });
 
  if (!response.ok) {
       throw new Error('Error when updating a task');
  }
 }
 
 // Функция удаления задачи.
 async function deleteTask(taskId) {
  const response = await fetch(`/api/tasks/delete?id=${taskId}`, {
       method: 'DELETE',
  });
 
  if (!response.ok) {
       throw new Error('Error when deleting a task');
  }
 }
 
 // Функция обновления списка задач.
 async function refreshTaskList() {
  const statusFilter = document.getElementById('status-filter').value;
  const sortOrder = document.getElementById('sort-filter').value;
  const response = await fetch(`/api/tasks?status=${statusFilter}&sort=${sortOrder}`);
  const tasks = await response.json();
  const taskList = document.getElementById('task-list');
   
  taskList.innerHTML = '';
   
  tasks.forEach(task => {
       const taskItem = createTaskItem(task);
       taskItem.querySelector('.status-select').value = task.status;
       taskList.appendChild(taskItem);
  });
 }
 
 // Функция создания элемента задачи.
 function createTaskItem(task) {
  const taskItem = document.createElement('li');
  taskItem.classList.add('task-item');
  taskItem.dataset.taskId = task.id;
   
  taskItem.innerHTML = `
       <span class="task-text">${task.text}</span>
       <span class="task-date">${new Date(task.date).toLocaleString()}</span>
       <input type="text" class="edit-input" style="display: none;">
       <select class="status-select">
           <option value="в процессе" ${task.status === 'в процессе' ? 'selected' : ''}>В процессе</option>
           <option value="завершено" ${task.status === 'завершено' ? 'selected' : ''}>Завершено</option>
       </select>
       <button class="edit-btn">Редактировать</button>
       <button class="delete-btn">Удалить</button>
  `;
   
  return taskItem;
 }
 