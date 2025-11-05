#!/bin/bash
services_directories=( services/*/ )
services=()

# echo "service dirs: ${services_directories[@]}"
          
for service_dir in "${services_directories[@]}"; do
  # get all directories in services folder and strip away end '/' and 'services/' prefix
  service=$(echo "${service_dir}" | sed -E 's/^services\///' | sed -E 's/\/$//')
  services+=("$service")
done

changed_services=()
ignored_folders_and_files=( "infrastructure" "README.md" )

for service in "${services[@]}"; do
  # Check if there are changes in the service directory
  # TODO: exclude certain directories from being considered if they are the only change
  service_in_changes=$(git diff --name-only HEAD~1 HEAD | grep -E "^services/${service}/")
  if [ -n "$service_in_changes" ]; then
    filtered_changes=$service_in_changes
    for ignored in "${ignored_folders_and_files[@]}"; do
      filtered_changes=$(echo "${filtered_changes}" | grep -vE "^services/${service}/${ignored}")
    done
    if [ -n "$filtered_changes" ]; then
      changed_services+=("$service")
    fi
  fi
done

echo "${changed_services[@]}"