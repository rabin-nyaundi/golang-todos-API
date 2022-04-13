echo "You are publishing your changes to github."

git add .

git status

echo "Please enter commit message."

read message 

branch=$(git symbolic-ref --short HEAD)

echo "You are on the" $branch " branch"

git commit -m "${message}"

if [-n "$(git status - porcelain)" ];

then
    echo "Its clean, all changes committed"
else
    git status

    echo "Pushing code to github!!!"
    
    git push -u origin $branch
fi


