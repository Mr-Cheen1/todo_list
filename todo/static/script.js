// Обработчик события DOMContentLoaded.
document.addEventListener('DOMContentLoaded', async function() {
  await refreshTaskList();
});

// Обработчик отправки формы создания задачи.
document.getElementById('task-form').addEventListener('submit', async function(e) {
  e.preventDefault();
  const taskInput = document.getElementById('task-input');
  const taskText = taskInput.value.trim();
  const expectedDateInput = document.getElementById('expected-date-input');
  const expectedDate = expectedDateInput.value;

  if (taskText === '') {
    alert('Введите текст задачи');
    return;
  }

  if (expectedDate === '') {
    alert('Выберите планируемую дату завершения задачи');
    return;
  }

  if (taskText.length > 255) {
    alert('Текст задачи не может превышать 255 символов');
    return;
  }

  const currentDate = new Date();
  currentDate.setUTCHours(0, 0, 0, 0); // Установка времени в полночь по UTC

  if (new Date(expectedDate + 'T00:00:00Z') < currentDate) {
    alert('Планируемая дата завершения не может быть раньше текущей даты');
    return;
  }

  const task = {
    text: taskText,
    createdDate: currentDate.toISOString().slice(0, 10),
    expectedDate: new Date(expectedDate + 'T00:00:00Z').toISOString().slice(0, 10),
    status: 0
  };

  try {
    await createTask(task);
    taskInput.value = '';
    expectedDateInput.value = '';
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
    const expectedDateInput = taskItem.querySelector('.expected-date-input');
    const statusSelect = taskItem.querySelector('.status-select');

    if (editInput.style.display === 'none') {
      editInput.style.display = 'inline';
      expectedDateInput.style.display = 'inline';
      editInput.value = taskText.textContent;
      expectedDateInput.value = taskItem.querySelector('.task-expected-date').textContent;
      taskText.style.display = 'none';
      taskItem.querySelector('.task-expected-date').style.display = 'none';
      statusSelect.style.display = 'none';
      e.target.textContent = 'Сохранить';
    } else {
      const updatedTask = {
        id: parseInt(taskId),
        text: editInput.value.trim(),
        createdDate: new Date(taskItem.querySelector('.task-created-date').textContent).toISOString().slice(0, 10),
        expectedDate: expectedDateInput.value? new Date(expectedDateInput.value).toISOString().slice(0, 10) : null,
        status: parseInt(statusSelect.value)
      };

      // Валидация полей задачи.
      if (updatedTask.text === '') {
        alert('Текст задачи не может быть пустым');
        return;
      }

      if (expectedDateInput.value === '') {
        alert('Выберите планируемую дату завершения задачи');
        return;
      }

      if (updatedTask.expectedDate && new Date(updatedTask.expectedDate) < new Date(updatedTask.createdDate)) {
        alert('Планируемая дата завершения не может быть раньше даты создания задачи');
        return;
      }

      if (updatedTask.text.length > 255) {
        alert('Текст задачи не может превышать 255 символов');
        return;
      }

      try {
        await updateTask(updatedTask);
        taskText.textContent = editInput.value;
        taskItem.querySelector('.task-expected-date').textContent = expectedDateInput.value;
        editInput.style.display = 'none';
        expectedDateInput.style.display = 'none';
        taskText.style.display = 'inline';
        taskItem.querySelector('.task-expected-date').style.display = 'inline';
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
    const taskId = parseInt(taskItem.dataset.taskId);
    const taskText = taskItem.querySelector('.task-text').textContent;
    const statusSelect = taskItem.querySelector('.status-select');

    const updatedTask = {
      id: taskId,
      text: taskText,
      createdDate: new Date(taskItem.querySelector('.task-created-date').textContent).toISOString().slice(0, 10),
      expectedDate: new Date(taskItem.querySelector('.task-expected-date').textContent).toISOString().slice(0, 10),
      status: parseInt(statusSelect.value)
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
  console.log('Updating task:', task);
  const response = await fetch(`/api/tasks/update?id=${parseInt(task.id)}`, {
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
  const response = await fetch(`/api/tasks?status=${statusFilter}&sort=${sortOrder}&sortField=createdDate`);
  const tasks = await response.json();
  const taskList = document.getElementById('task-list');

  taskList.innerHTML = '';

  if (tasks && tasks.length > 0) {
    // Создание и добавление контейнера для заголовков
    const headers = document.createElement('div');
    headers.className = 'task-headers';
    headers.style.display = 'flex';
    headers.style.justifyContent = 'between';
    headers.style.width = '100%';

    const headerText = document.createElement('div');
    headerText.className = 'task-header';
    headerText.textContent = 'Наименование задачи';
    headerText.style.flex = '1';
    headerText.style.maxWidth = '198px'; 
    headerText.style.display = 'flex'; 
    headerText.style.alignItems = 'center'; 
    headerText.style.justifyContent = 'center';

    const headerCreated = document.createElement('div');
    headerCreated.className = 'task-header';
    headerCreated.textContent = 'Дата добавления';
    headerCreated.style.flex = '1';
    headerCreated.style.maxWidth = '97px'; 

    const headerExpected = document.createElement('div');
    headerExpected.className = 'task-header';
    headerExpected.textContent = 'Дата завершения';
    headerExpected.style.flex = '1';
    headerExpected.style.maxWidth = '102px'; 

    const headerStatus = document.createElement('div');
    headerStatus.className = 'task-header';
    headerStatus.textContent = 'Статус задачи';
    headerStatus.style.flex = '1';
    headerStatus.style.maxWidth = '130px';
    headerStatus.style.display = 'flex'; 
    headerStatus.style.alignItems = 'center'; 
    headerStatus.style.justifyContent = 'center';

    headers.appendChild(headerText);
    headers.appendChild(headerCreated);
    headers.appendChild(headerExpected);
    headers.appendChild(headerStatus);
    taskList.appendChild(headers);

    tasks.forEach(task => {
      const taskItem = createTaskItem(task);
      taskList.appendChild(taskItem);
    });
  } else {
    const emptyMessage = document.createElement('li');
    emptyMessage.textContent = 'Нет задач для отображения.';
    emptyMessage.style.textAlign = 'center';
    emptyMessage.style.fontStyle = 'italic';
    emptyMessage.style.color = 'gray';
    taskList.appendChild(emptyMessage);
  }
}

// Функция создания элемента задачи.
function createTaskItem(task) {
  const taskItem = document.createElement('li');
  taskItem.classList.add('task-item');
  taskItem.dataset.taskId = task.id;

  taskItem.innerHTML = `
    <div class="task-text">${task.text}</div>
    <div class="task-created-date">${task.createdDate}</div>
    <div class="task-expected-date">${task.expectedDate}</div>
    <input type="text" class="edit-input" style="display: none;">
    <input type="date" class="expected-date-input" style="display: none;">
    <select class="status-select">
      <option value="0" ${task.status === 0 ? 'selected' : ''}>В процессе</option>
      <option value="1" ${task.status === 1 ? 'selected' : ''}>Завершено</option>
      <option value="2" ${task.status === 2 ? 'selected' : ''}>Тестирование</option>
      <option value="3" ${task.status === 3 ? 'selected' : ''}>Возвращено</option>
    </select>
    <button class="edit-btn">Редактировать</button>
    <button class="delete-btn">Удалить</button>
  `;

  return taskItem;
}
