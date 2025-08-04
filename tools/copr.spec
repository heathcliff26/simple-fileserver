%global debug_package %{nil}

Name:           simple-fileserver
Version:        0
Release:        %autorelease
Summary:        Minimal http file server
%global package_id io.github.heathcliff26.%{name}

License:        Apache-2.0
URL:            https://github.com/heathcliff26/%{name}
Source:         %{url}/archive/refs/tags/v%{version}.tar.gz

BuildRequires: golang >= 1.24

%global _description %{expand:
A golang implementation of a minimal http file server.}

%description %{_description}

%prep
%autosetup -n %{name}-%{version} -p1

%build
export RELEASE_VERSION="%{version}-%{release}"
make build

%install
install -D -m 0755 bin/%{name} %{buildroot}%{_bindir}/%{name}
install -D -m 0644 tools/%{name}.service %{buildroot}%{_prefix}/lib/systemd/system/%{name}.service
install -D -m 0644 %{package_id}.metainfo.xml %{buildroot}/%{_datadir}/metainfo/%{package_id}.metainfo.xml

%post
systemctl daemon-reload

%preun
if [ $1 == 0 ]; then #uninstall
  systemctl unmask %{name}.service
  systemctl stop %{name}.service
  systemctl disable %{name}.service
  echo "Clean up %{name} service"
fi

%postun
if [ $1 == 0 ]; then #uninstall
  systemctl daemon-reload
  systemctl reset-failed
fi

%files
%license LICENSE
%doc README.md
%{_bindir}/%{name}
%{_prefix}/lib/systemd/system/%{name}.service
%{_datadir}/metainfo/%{package_id}.metainfo.xml

%changelog
%autochangelog
