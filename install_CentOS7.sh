#!/bin/tcsh -f
echo "Installing rvlink protocol handler for XDG"

set desktopInstall = "`which desktop-file-install`"

if ($status == 0) then
    set xdgDir = $HOME/.local/share/applications
    if ( ! -e $xdgDir ) then
        mkdir -p $xdgDir
    endif

    if ( -e $xdgDir) then
        set xdgFile = $xdgDir/dilink.desktop
        rm -f $xdgFile
        cat > $xdgFile << EOF
[Desktop Entry]
Name=DiLink
Type=Application
Exec=/lustre/INHouse/CentOS/bin/dilink %U
Terminal=false
Categories=AudioVideo;Viewer;Player;
MimeType=x-scheme-handler/dilink;
NoDisplay=true
EOF
        echo "Successfully created ${xdgFile}"
    else
        echo "WARNING: can't find or create XDG directory: ${xdgDir}:  skipping XDG url-handler registration."
    endif

    desktop-file-install $xdgFile --rebuild-mime-info-cache --dir=$xdgDir
else
    echo "WARNING: desktop-file-install not found: skipping xdg-handler registration."
endif

echo "Done."
